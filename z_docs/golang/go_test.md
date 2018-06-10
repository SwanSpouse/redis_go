

### Golang Test

#### ginkgo

Ginkgo Homepage: http://onsi.github.io/ginkgo/

Ginkgo is a BDD-style Go testing framework built to help you efficiently write expressive and comprehensive tests. It is best paired with the Gomega matcher library but is designed to be matcher-agnostic.

* Ginkgo makes extensive use of closures to allow you to build descriptive test suites.

* You should make use of Describe and Context containers to expressively organize the behavior of your code.

Ginkgo has support for running specs in parallel. It does this by spawning separate go test processes and serving specs to each process off of a shared queue. This is important for a BDD test framework, as the shared context of the closures does not parallelize well in-process.
