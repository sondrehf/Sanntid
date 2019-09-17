# Reasons for concurrency and parallelism


To complete this exercise you will have to use git. Create one or several commits that adds answers to the following questions and push it to your groups repository to complete the task.

When answering the questions, remember to use all the resources at your disposal. Asking the internet isn't a form of "cheating", it's a way of learning.

 ### What is concurrency? What is parallelism? What's the difference?
 > Concurrency – processes executed at the same time, but it is sufficient for the processes to appear to be executed simultaneously. Parallellism - processes run in parallell. 
 
 ### Why have machines become increasingly multicore in the past decade?
 > Increased processing speed, more energy efficient. 

 
 ### What kinds of problems motivates the need for concurrent execution?
 (Or phrased differently: What problems do concurrency help in solving?)
 > Separate processes and independent processes,  which should be executed seemingly simultaneously. It can advance without waiting for all other computations to complete. 

 
 ### Does creating concurrent programs make the programmer's life easier? Harder? Maybe both?
 (Come back to this after you have worked on part 4 of this exercise)
 > Both, it introduces a new set of problems/issues, but also make it easier for some classes of problems. Multithreading makes solving many processes easier, but can make the debugging process more difficult. 

 
 ### What are the differences between processes, threads, green threads, and coroutines?
 > Process – OS-managed with own address space. 
 > Threads – OS-managed within same address space. 
 > Green threads – User-managed 
 > Coroutines – done sequentially, not in parallel as with the others. 
 
 ### Which one of these do `pthread_create()` (C/POSIX), `threading.Thread()` (Python), `go` (Go) create?
 > The first two creates a thread, go create a couroutine. 
 
 ### How does pythons Global Interpreter Lock (GIL) influence the way a python Thread behaves?
 > It prevents multiple threads to access the Python bytecode at the same time. 
 
 ### With this in mind: What is the workaround for the GIL (Hint: it's another module)?
 > The multiprocessing package
 
 ### What does `func GOMAXPROCS(n int) int` change? 
 > Increases the amount of allocated operating system threads in a Go program 

