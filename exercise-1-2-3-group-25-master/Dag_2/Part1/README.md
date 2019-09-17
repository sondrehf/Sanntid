# Mutex and Channel basics

### What is an atomic operation?
> Concurrent operations with no interrupts, completed in a single step. The different operations run completely independently of the other processes. For instance fetch and add.

### What is a semaphore?
> To avoid race conditions, a semaphore is a variable to controll access to a common resource between threads. 

### What is a mutex? (MUTualEXclusion)
> Used in concurrency to prevent race conditions, so that two threads can not access a critical selection(see question 5) at the same time. 

### What is the difference between a mutex and a binary semaphore?
> Mutex - only released from the thread that aquired it. 
> Binary semaphore - Can be signaled from any other thread. 

### What is a critical section?
> Program part where a thread accesses shared resourses. Therefore, a synchronization method must be implemented to avoid race conditions.  

### What is the difference between race conditions and data races?
 > Race condition - timing fault that leades to weird/wrong behaviour. 
 > Data races - two different threads writing to the same data location at almost the same time. At least one is writing

### List some advantages of using message passing over lock-based synchronization primitives.
> - Easier
> - Less programmer errors
> - With more threads, not an exponentiall increase in complexity
> - Can be event driven

### List some advantages of using lock-based synchronization primitives over message passing.
> - Better performance
> - No need to allocate message objects
