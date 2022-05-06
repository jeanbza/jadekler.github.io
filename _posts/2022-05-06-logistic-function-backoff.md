---
layout: post
title:  "Graceful degradation with the logistic function"
date:   2022-05-05 15:55:23 -0600
categories: 
---

# Graceful degradation with the logistic function

I recently worked on a server throttling feature in one of our build stack's busiest binaries. We serve files from this binary, which is deployed as tens of thousands of tasks, which cumulatively serve millions of QPS.

Sometimes, one of these tasks gets a large memory spike. The cause for this is that the task is asked to hold a single file, and before it gets replicated, the file is needed by tens of thousands of peers, who all bombard the task with RPCs. Handling all these requests causes a huge spike in memory, and the task falls over.

Eventually replication catches up, and there are enough peers to spread the load. But we'd like to fail more gracefully than just OOM => death. We'd like to monitor the current and allocated memory, and gradually reject requests (throttle) when current memory exceeds allocated memory.

Technical note: the binary gets deployed to a container that hosts several processes, all of which share the same memory. Each task has some allocated memory, but can exceed their allocation. Doing so chews into other processes' memory. One process going far above allocation and OOMing means that all processes OOM.

# Abrupt degradation

A simple way to gracefully degrade is to reject requests above the allocated amount. We can model that with a simple step function. For the sake of example, let's imagine that our allocated memory is `3GiB = 3221225472 bytes`. We don't want to hit 3GiB exactly, since that's roughly where we'll OOM, so let's start throttling a bit before that: at `3000000000 bytes = 3e9 bytes`.

$$
f(x) = \left\{
        \begin{array}{ll}
            1 & \quad x \ge 3e9 \\
            0 & \quad x < 3e9
        \end{array}
    \right.
$$

Here, 0 means "don't reject", and 1 means "reject".

But, this is inefficient:

- We're not using all our available memory. In a resource constrained environment, or when we're highly scaled, we really want to squeeze every bit of memory that we can.
- With the more flexible memory model mentioned above, in which we can steal some memory from other tasks, we want to allow our process to _exceed_ its allocated memory. The more it exceeds its allocation, the more we should throttle, until we reach an unacceptably high excess, at which point we should throttle all requests.

# Gradual degradation

It's starting to sound like we need a linear function, not a stepwise function: something that rejects more and more requests the more memory we're using. We now need a range to operate our throttler within: at the bottom of the range, we reject no requests; at the top, all requests.

Using the rigid memory model above, we'll define our range as `[3e9 bytes, 3221225472 bytes]`. In the more flexible memory model our production code operates under, it's more like `[max_bytes, max_bytes+allowed_theft]`.

So, let's build a linear function for this. Note that the values that we want from our linear function are `[0.0, 1.0]`. As above, 0 means "don't reject", and 1 means "reject". Any value between that represents that chance that a request will be rejected. That is, we'll compare our the result of our linear function against a number taken randomly from a uniform distribution of `[0.0, 1.0]`.

To build this linear function, let's start with what we know:

- Linear functions look like $f(x)=a \cdot x + b$
- $f(3e9)=0$
- $f(3221225472)=1$

We can use this to solve the equation:

$$
\begin{split}
0 = a \cdot 3e9 + b\\
-b = 3e9 \cdot a\\
b = 3e9 \cdot a\\
b = -3e9 \cdot a
\end{split}
\quad\quad
\begin{split}
1 = a \cdot 3221225472 + b\\
1 - b = a \cdot 3221225472\\
-b = a \cdot 3221225472 - 1\\
b = -a \cdot 3221225472 + 1
\end{split}\\ \text{ } \\ \text{ } \\
-3e9 \cdot a = -a \cdot 3221225472 + 1 \\
-3e9 \cdot a + a \cdot 3221225472 = 1 \\
a \cdot (-3e9 + 3221225472) = 1 \\
a = 1/221225472
$$

Now that we know `a`, let's use that and either of the two partial solutions above to find `b`:

$$
1 = a \cdot 3221225472 + b\\
1 = 3221225472/221225472 + b\\
b = 1 - 3221225472/221225472\\
b = -1953125/144027
$$

Our stepwise function accordingly gets an embedded linear function:

$$
f(x) = \left\{
        \begin{array}{ll}
            x/221225472 - 1953125/144027  & \quad x \ge 3e9 \\
            0 & \quad x < 3e9
        \end{array}
    \right.
$$

![](/assets/linear.png)

We can verify that this works by plugging in our original numbers:

$$
\begin{split}
f(x) = x/221225472 - 1953125/144027\\
f(x) = 3e9/221225472 - 1953125/144027\\
f(x) = 0
\end{split}
\quad\quad
\begin{split}
f(x) = x/221225472 - 1953125/144027\\
f(x) = 3221225472/221225472 - 1953125/144027\\
f(x) = 1
\end{split}\\
$$

# Graceful degradation

This is a lot nicer, but let's consider an _even more_ resource constrained environment, in which memory usage consistently hovers near the threshold. Let's also consider a binary that handles a wide range of requests, whose memory requirements vary significantly.

In such an environment, we'd like to throttle as few requests as possible: especially when we're on the lower end of our memory throttle range. The reasoning here is that with a varied workload, we're not sure that a linear increase in requests will result in a linear increase in memory. It may be that for a period of time, we're running into the excess space, but the incoming requests take dramatically less memory than prior requests: we want to throttle as few of these as possible. This lets us take full advantage of our overflow range. But, if we do start creeping further into that range, we want to quickly start throttling more aggressively.

The [logistic function](https://en.wikipedia.org/wiki/Logistic_function) is perfect for that. Here's the shape of the logistic function:

![](/assets/logistic.png)

Its S-shaped curve allows more requests through when we're at the bottom of our range, and aggressively throttles at the end of our range.

The equation for the logistic function is as follows,

$$
\begin{align*}
f(x) = \dfrac{L}{1 + e^{-k(x-x_0)}}
\end{align*}
$$

Where,

- $x_0$ is the x value of the sigmoid's midpoint
- L is the curve's maximum value
- k is the logistic growth rate or steepness of the curve

Let's adapt this to our problem:

L is the easiest: we want the maximum value to be 1 (we want our range to be `[0.0, 1.0]`), So, `L=1`.

$x_0$ is fairly straightforward: the midpont should be the midpoint between the start and end of our range. So,

$$
\begin{align}
x_0 = 3221225472-\left(\dfrac{3221225472-3e9}{2}\right)\\
x_0 = 3110612736
\end{align}
$$

Now we have,

$$
\begin{align*}
f(x) = \dfrac{L}{1 + e^{-k(x-x_0)}}\\
f(x) = \dfrac{1}{1 + e^{-k(x-3110612736)}}
\end{align*}
$$

k is the hardest. Let's start by solving for k in the equation above:

$$
\begin{align*}
f(x) = \dfrac{1}{1 + e^{-k\left(x-3110612736\right)}}\\
f(x)\left(1 + e^{-k(x-3110612736)}\right) = 1\\
1 + e^{-k(x-3110612736)} = \dfrac{1}{f(x)}\\
e^{-k(x-3110612736)} = \dfrac{1}{f(x)}-1\\
-k(x-3110612736) = ln\left(\dfrac{1}{f(x)}-1\right)\\
k(x-3110612736) = -ln\left(\dfrac{1}{f(x)}-1\right)\\
k = \dfrac{-ln\left(\dfrac{1}{f(x)}-1\right)}{x-3110612736}\\
\end{align*}
$$

Now we return to what we know about how this curve _should_ behave:

- $f(3e9)=0$
- $f(3221225472)=1$

Unfortunately, using either of these results in an unsolvable equation:

$$
\begin{split}
k = \dfrac{-ln\left(\dfrac{1}{f(x)}-1\right)}{x-3110612736}\\
k = \dfrac{-ln\left(\dfrac{1}{0}-1\right)}{3e9-3110612736}\\
\text{NaN: can't divide by 0}
\end{split}
\quad\quad
\begin{split}
k = \dfrac{-ln\left(\dfrac{1}{f(x)}-1\right)}{x-3110612736}\\
k = \dfrac{-ln\left(\dfrac{1}{1}-1\right)}{3221225472-3110612736}\\
k = \dfrac{-ln(0)}{3221225472-3110612736}\\
\text{NaN: natural log of 0 is undefined}
\end{split}
$$

So, that's a bummer. But it makes sense: the logistic function is asymptotic, with asymptotes 0 and 1: it will never actually reach those values!

So, let's estimate k by choosing a value close to the asymptotes: either .01 for the lower bound, or .99 for the upper bound. It doesn't matter which one we do, as the curve is reflected around the midpoint. So, let's use the upper:

$$
k = \dfrac{-ln\left(\dfrac{1}{f(x)}-1\right)}{x-3110612736}\\
k = \dfrac{-ln\left(\dfrac{1}{.99}-1\right)}{3221225472-3110612736}\\
k = \dfrac{-ln(0.0101010101)}{110612736}\\
k = \dfrac{4.59511985023}{110612736}\\
k = \dfrac{4.595119}{110612736}\\
k = .0000000415424043
$$

Great! Let's put it all together:

$$
\begin{align*}
f(x) = \dfrac{L}{1 + e^{-k(x-x_0)}}\\
f(x) = \dfrac{1}{1 + e^{-.0000000415424043(x-3110612736)}}
\end{align*}
$$

![](/assets/logistic_real.png)

We can verify that this works by plugging in our original numbers:

$$
\begin{split}
f(x) = \dfrac{1}{1 + e^{-.0000000415424043(x-3110612736)}}\\
f(x) = \dfrac{1}{1 + e^{-.0000000415424043(3e9-3110612736)}}\\
f(x) = 0.01000000841
\end{split}
\quad\quad
\begin{split}
f(x) = \dfrac{1}{1 + e^{-.0000000415424043(x-3110612736)}}\\
f(x) = \dfrac{1}{1 + e^{-.0000000415424043(3221225472-3110612736)}}\\
f(x) = 0.98999999158
\end{split}\\
$$

Since we're approximating values and will never reach 0 or 1, it's helpful to continue using the stepwise function to guarantee no throttling when we're below our threshold, and to always throttle when we're above our allowable range:

$$
f(x) = \left\{
        \begin{array}{ll}
            1 & \quad x \ge 3221225472 \\
            \dfrac{1}{1 + e^{-.0000000415424043(x-3110612736)}}  & \quad 3e9 \le x < 3221225472  \\
            0 & \quad x < 3e9
        \end{array}
    \right.
$$

# Conclusion

At tremendous scale, it's important to eke every last bit of memory from servers. It's also important to be able to gracefully degrade during memory spikes, to avoid out-of-memory crashes. The logistic function is an excellent function for deciding whether to throttle requests, which strikes a good balance between the competing priorities of using all available memory and avoiding out-of-memory crashes.

<script type="text/javascript" async
  src="/assets/MathJax-2.7.9/MathJax.js?config=TeX-AMS-MML_HTMLorMML">
</script>
<script type="text/x-mathjax-config">
    MathJax.Hub.Config({
      tex2jax: {
        inlineMath: [ ['$','$'], ["\\(","\\)"] ],
        processEscapes: true
      },
      extensions: [
        "MathMenu.js",
        "MathZoom.js",
        "AssistiveMML.js"
      ],
      jax: ["input/TeX", "output/CommonHTML"],
      TeX: {
        extensions: [
          "AMSmath.js",
          "AMSsymbols.js",
          "noErrors.js",
          "noUndefined.js",
        ]
      }
    });
  </script>