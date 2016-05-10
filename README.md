# pcg
Go implementation of Melissa O'Neill's excellent PCG pseudorandom number generator, which is 
well-studied, excellent, and fast both to create and in execution.

````
  Performance on a MacBook Pro:

  $ go test -v -bench=.
  === RUN   TestSanity32
  --- PASS: TestSanity32 (0.00s)
  === RUN   TestSum32
  --- PASS: TestSum32 (0.00s)
  === RUN   TestAdvance32
  --- PASS: TestAdvance32 (0.00s)
  === RUN   TestRetreat32
  --- PASS: TestRetreat32 (0.00s)
  === RUN   TestSanity64
  --- PASS: TestSanity64 (0.00s)
  === RUN   TestSum64
  --- PASS: TestSum64 (0.00s)
  === RUN   TestAdvance64
  --- PASS: TestAdvance64 (0.00s)
  === RUN   TestRetreat64
  --- PASS: TestRetreat64 (0.00s)
  === RUN   ExampleReport32
  --- PASS: ExampleReport32 (0.00s)
  === RUN   ExampleReport64
  --- PASS: ExampleReport64 (0.00s)
  BenchmarkNew32-8      2000000000               1.09 ns/op
  BenchmarkRandom32-8   1000000000               2.49 ns/op
  BenchmarkBounded32-8  200000000                9.75 ns/op
  BenchmarkNew64-8      200000000                6.89 ns/op
  BenchmarkRandom64-8   200000000                7.58 ns/op
  BenchmarkBounded64-8  50000000                25.5 ns/op
````
