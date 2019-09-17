---
layout: post
title:  "Code Reviews: Gerrit vs GitHub"
date:   2020-08-21 15:55:23 -0600
categories: 
---

I've now used Gerrit and GitHub for several years each, and I've developed
Opinions about these code review systems that I wanted to write down.

Note: this is written in August of 2019. Hopefully, if you're reading this in
the future, both systems have improved. (but especially GitHub)

# Goals

To start, I've bucketed the types of things I want out of a code review system
in different buckets, and rathed them based on a 5 point scale:

- Very bad: ❌❌
- Bad: ❌
- Meh: 🤷
- Good: ✅
- Very good: ✅✅

TODO: create table for the remainder

**As devops**

- Easy (or zero work) to administer the system.
- Ability to moderate the system.
  - Easy to add/remove people, groups of people.
  - Different levels of approvals (some people can thumbs up, subset can
  submit).
- Authentication is provided elegantly.

**As a contributor**

- Easy to submit a change.

**As a reviewer**

- 

**As either contributor or reviewer**

- A dashboard to see my pending changes and things that are waiting for my
review.
- Easy to parse reviewer's comments.
  - Easy to understand which line(s) a comment relates to.
  - Ability to close/resolve comments for things that I've finished.
  - Ability to see unclosed/unresolved comments.
  - Large CLs with many rounds of review should still be fairly easily
  comprehensible.
- Ability to assign multiple people to a review.
- Ability to require multiple people's thumbs-up for a review.