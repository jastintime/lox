# Lox
The main branch holds the golang implementation of lox and the C branch will have the C implementation.

# TODO
Pass all tests in the crafting interpreters test suite.

| Test | Current Status |
| ---- | -------------- |
| chap08\_statements | 69/69 |
| chap09\_control | 95/96 | 
| chap10\_functions | 136/139 |
| chap11\_resolving | 144/147 |
| chap12\_classes | 204/207 |
| chap13\_inheritance | 236/239 |

Note, test are cumulative, a test not passing
for chap09 will also be counted in the chap10 test

Currently stuck passing the last three test, these all have to do with scoping in for loops,
the culprit seems to be that i.Locals is being modified and the one at the correct distance
is being replaced. I've spend far too long trying to figure out what line is wrong, a
pull request would be greatly appreciated.
