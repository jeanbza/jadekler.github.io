---
layout: post
title:  "Binary Search Tree Runtime Complexities"
date:   2019-11-23 15:55:23 -0600
categories: 
---

# Introduction

I've recently inundated myself with interview preparation. Along the way, I
thought a lot about how to intuit runtime complexities for various algorithms.
I thought that it might be nice to cement it all with an article - both for the
sake of others, as well as for my sake. This is the first article in that vein,
strictly dealing with binary search trees.

# Terminology

Note: This article frequently uses the abbreviation BST to describe binary
search trees.

Here are properties of trees that this article deals with:

- Number of nodes in a tree is typically denoted `n`.
- Height of a tree is typically denoted `k`.
- [Complete](https://en.wikipedia.org/wiki/Binary_tree): every level except
  possibly the last is completely filled.
- [Full](https://en.wikipedia.org/wiki/Binary_tree): every node has 0 or 2
  children. This is a subset of complete: every full tree is a complete tree.
- [Balanced](https://en.wikipedia.org/wiki/Binary_tree): the height of the left
  and right subtrees of every node differ by at most 1. Using the root node as
  our focus, we can intuit that any balanced tree has fairly uniform height.
- [Binary Search Tree](https://en.wikipedia.org/wiki/Binary_search_tree):
  A tree in which each element to the left of a node is guaranteed to be less,
  and to the right guaranteed to be greater. _BSTs may not be balanced_.
- [Red-black Tree](https://en.wikipedia.org/wiki/Red%E2%80%93black_tree): A BST
  that is balanced.
- [AVL Tree](https://en.wikipedia.org/wiki/AVL_tree): A BST that is balanced (in
  a different way than red-black trees).
- [Branching factor](https://en.wikipedia.org/wiki/Branching_factor): The number
  of children at each node.

# Searching a BST

Everyone knows that searching a binary search tree has runtime complexity
`O(logn)`, right? ...right? Let's take a second to tease apart some questions
from that assertion to see if we really understand what we mean when we say that:

- Q: What is the [base](https://en.wikipedia.org/wiki/Logarithm) of the log?

  A: Typically when we talk about the base of a log, we're talking about 2 or 10.  In the case of a BST case it's 2 (we'll dive into why that is shortly).

- Q: What does `n` represent?

  A: `n` _usually_ means "number of elements". In this case, `n` does mean that: or, another way to put that is "number of nodes in the tree".

So, expanding `O(logn)`, we have: `O(log2(<# nodes in tree>))`.

Well, searching a BST is not strictly `O(log2(n))`: it depends on whether the
BST is balanced or not. An unbalanced BST may be `O(n)`. Let's discover why by
exploring how the `log2` comes about.

# Logarithms and exponents

Let's look at the following _full_ binary tree:

![Full Binary Tree](/assets/simple_complete.png)

We can tell the following about this tree:

- It has height `k=4`.
- It has nodes `n=15`.

The number of nodes `n` can also be expressed as,

$$
\begin{align*}
n=2^k-1 \\
n=2^4-1 \\
n=16-1 \\
n=15
\end{align*}
$$

This is easily observable: each level has 2 times the number of nodes as the
last, since each node has 2 children (except the leaf nodes, which have 0).

What if we knew the amount of nodes, but not the height? We can rework the
formula we just came up with to give us height using number of nodes,

$$
\begin{align*}
n=2^k-1 \\
n+1=2^k \\
log2(n+1)=log2(2^k) \\
log2(n+1)=k*log2(2) \\
log2(n+1)=k \\
k=log2(n+1)
\end{align*}
$$

# Back to binary search trees

How can we apply this knowledge? Consider again our original topic of a binary
search tree: one property of a BST is that a search operation requires no
backtracking. That is: the path to a node always going to go at most to a leaf
node - it never reaches a leaf and then has to backtrack and try a different
choice at a former node.

If the BST is unbalanced, the worst case time complexity to search for a value
is `O(n)`. We can easily show that this is the case with the following tree:

![Left BST](/assets/left_bst.png)

This BST is heavily weighted to the right. To search for the value 109, we'd have to
look at each of the `n=15` elements. That is: the worst case time complexity is
`O(n)`.

What if this BST were balanced, instead?

![Balanced BST](/assets/balanced_bst.png)

Now, to search for _any_ value in the BST, the maximum depth we'd need to
traverse is `k`. We know that `k=log2(n)`, so we can say that the big-O runtime
complexity in terms of `n` is `O(log2(n))`.

# An aside on branching factor

All of our details above deal with a binary tree - a tree in which each node has
at most 2 children. What if each node had 9 children? The maximum number of
children a node has in a tree is called the [branching factor](https://en.wikipedia.org/wiki/Branching_factor),
and it determines the bases of the logarithms and exponents in all the
calculations above.

So, if the branching factor for a tree is 9, and it is complete, then it has
`9^k-1` nodes (and so on).

# In conclusion

In conclusion: binary search trees don't guarantee `O(log2(n))` search:
_balanced_ BSTs do. Hopefully the above provides the intuition for that.

In a future article we'll look at [Heap](https://en.wikipedia.org/wiki/Heap_(data_structure))
runtime complexity, and how memoization affects that. In some other future
article we'll look at backtracking and permutation runtime complexities.

<script type="text/javascript" async
  src="/assets/mathjax-2.7.5/MathJax.js?config=TeX-MML-AM_CHTML">
</script>
<script type="text/x-mathjax-config">
    MathJax.Hub.Config({
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