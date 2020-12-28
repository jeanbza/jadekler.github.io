---
layout: post
title:  "Permutations, memoization, and their runtime complexity"
date:   2020-10-18 15:55:23 -0600
categories: 
---

# Introduction

I've frequently struggled with analysing the runtime complexity of a permutation
based algorithm. I recently came across a really great interview question that
involves permutations and memoization (caching the return value of repeated
calls to a function). It helped make permutation runtime complexity "click" for
me, and the memoization runtime complexity parts were really cool. 🤓

So, I thought I'd write about it. Let's dive in and talk about permutations,
memoization, and runtime complexity!

# Permutations

Let's start by considering an obvious permutation problem:

```
Given an input set of unique characters, how many unique strings of length n can
be made? Each character in the input can only be used once.

Ex: given ABC, we can make,

ABC
ACB
BAC
BCA
CAB
CBA

So, Permutations("ABC") should return 6.
```

This is a fairly straightforward permutation problem. Sometimes permutation
problems are marginally harder to spot, though, like
[knight dialer](https://hackernoon.com/google-interview-questions-deconstructed-the-knights-dialer-f780d516f029).
Incidentally, knight dialer inspired this post, and the linked article has a
great write up! This article will be very similar in spirit to it.

So, back to our interview question: how many permutations can we make? Well, we
could generate all the permutations and count them. That's fairly
straightforward:

```go
func main() { fmt.Println(len(permutations([]rune("ABC"), 0))) }
func permutations(s []rune, start int) (combinations []string) {
  if start == len(s) {
	return []string{string(s)}
  }
  for i := start; i < len(s); i++ {
    s[start], s[i] = s[i], s[start] // swap
    combinations = append(combinations, permutations(s, start+1)...)
    s[start], s[i] = s[i], s[start] // undo the swap
  }
  return
}
```

[This works](https://play.golang.org/p/XgGprQQA1cf)! Let's analyse its runtime
complexity:

<div class="leftie">
$$
\begin{align}
\text{The first, depth 0 call iterates over $n$ items in $s$.} \\
\text{Each of those $n$ iterations spawns a depth=1 recursive function that loops over $n-1$ items.} \\
\text{Each of those $n-1$ iterations spawns a depth=2 recursive function that loops over $n-2$ items.} \\
\text{And so on, giving us a pattern of,} \\
n \cdot (n-1) \cdot (n-2) \cdot ... \cdot 1 \\
\sum_{i=1}^{n} i = n! \\
\end{align}
$$
</div>

So, the runtime complexity is `O(n!)`. And, we can infer the space complexity
(of our stack usage performing recursion) to be `O(n)`: it's not multi-threaded,
so only one recursion can happen at a time, and the max depth it can descend is
`n`.

----------
https://www.cs.sfu.ca/~ggbaker/zju/math/perm-comb-more.html#rep-comb
https://sites.math.northwestern.edu/~mlerma/courses/cs310-05s/notes/dm-gcomb
https://en.wikipedia.org/wiki/Binomial_coefficient

binomial coefficient is combination, order does not matter, no repetition
exponent is permutation, order matters, repetition allows

$$
\begin{align}
n=\sum\limits_{j=0}^{k-1} b^j \\
n=b^{k-1}+b^{k-2}+...+1 && \text{Expanding the sum.} \\
n \cdot b=(b^{k-1}+b^{k-2}+...+1) \cdot b && \text{Multiply by b.} \\
n \cdot b=b^k+b^{k-1}+...+b && \text{Distribute the b into the sum.} \\
n \cdot b + 1=b^k+b^{k-1}+...+b+1 && \text{Add 1 to each side.} \\
n \cdot b + 1=b^k + n && \text{Note that we can use n for the right side.} \\
n \cdot b - n=b^k - 1 && \text{Move n left; 1 right.} \\
n \cdot (b - 1)=b^k - 1 && \text{Factor out n.} \\
n = \frac{b^k - 1}{b - 1} && \text{Divide by n-1.} \\
\end{align}
$$

<script type="text/javascript" id="MathJax-script" async src="https://cdn.jsdelivr.net/npm/mathjax@3/es5/tex-chtml.js"></script>
<script type="text/javascript">
window.MathJax = {
  tex: {
    packages: ['base', 'ams', 'textmacros']
  },
  loader: {
    load: ['ui/menu', '[tex]/ams', '[tex]/textmacros']
  },
  tex: {packages: {'[+]': ['textmacros']}},
};
</script>

<style>
.leftie .MathJax * {
  text-align: left !important;
}
</style>
