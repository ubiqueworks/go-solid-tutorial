# SOLID Go demonstration

## Scenario
We want to build a simple application to manage a user collection. The service will expose the following functionalities:
* Create user
* Delete user
* List users

For business reasons, the **create** and **delete** functionalities will only be available to an administrative
interface via RPC while the **list** of users will be publicly available via HTTP 

## Implementations
STEP 0 - Flawed implementation

STEP 1 - Fix interfaces/dependency graph

STEP 2 - Fix models/domain per layer

STEP 3 - Fix configuration

## Guiding principles

*Two types are substitutable if they exhibit behaviour such that the caller is unable to tell the difference - Liskov principle*

*Require no more, promise no less – Jim Weirich*

*Clients should not be forced to depend on methods they do not use. – Robert C. Martin*

*High-level modules should not depend on low-level modules. Both should depend on abstractions. Abstractions should not depend on details. Details should depend on abstractions. – Robert C. Martin*

*In Go, your import graph must be acyclic. A failure to respect this acyclic requirement is grounds for a compilation failure, but more gravely represents a serious error in design.
 All things being equal the import graph of a well designed Go program should be a wide, and relatively flat, rather than tall and narrow.
 [...] The dependency inversion principle encourages you to push the responsibility for the specifics, as high as possible up the import graph, to your main package or top level handler, leaving the lower level code to deal with abstractions–interfaces. - Dave Cheney* 

## Reference documentation 

* [SOLID Go Design](https://dave.cheney.net/2016/08/20/solid-go-design)
* [Preemptive interface antipattern](https://medium.com/@cep21/preemptive-interface-anti-pattern-in-go-54c18ac0668a)
* [What “accept interfaces, return structs” means in Go](https://medium.com/@cep21/what-accept-interfaces-return-structs-means-in-go-2fe879e25ee8)

