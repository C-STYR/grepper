### How this all should work: 

 - a root path is entered on the command line

 - a recursive process gathers all filenames within the root

 - multiple routine start to search the contents of each file and send results down channel

 - a separate routine displays the results as they come in (working concurrently with the search routines)


#### `main.go`

 1. establish a master goroutine for waitgroups - it will receive from channels tied to all processes

 2. in a GR, start the recursive treebuilding in `buildTree.go`, and add to wg on each call

 3. once treebuilding waitgroup is finished, send signal to gatherFilenamesComplete channel (to master wg GR and to search GR)

 4. for each member of search party, add to wg and start GR for searching line by line
    - grab a task from the queue (each task is a filename)
    - create slice of results 
    - if slice is not empty, send slice down results channel
    - TODO: figure out how to stop all searchloops when treebuilding complete
    - as each GR completes, decrement wg

5. start GR to wait for completion of search loops. when all search loops have completed, send signal to master wg GR

6. 
