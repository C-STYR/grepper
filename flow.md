### How this all should work: 

 - a root path is entered on the command line

 - a recursive process gathers all filenames within the root

 - multiple routine start to search the contents of each file and send results down channel

 - a separate routine displays the results as they come in (working concurrently with the search routines)


#### `main.go`

 1. establish a master WG - it will wait for 3 main processes to complete

 2. in a GR, start the recursive treebuilding in `buildTree.go`

 3. once treebuilding is complete, close treebuilding channel and decrement master WG

 4. for each member of search party, start GR for searching line by line
    - grab a task from the queue (each task is a filename)
    - create slice for each filename and populate with search results
    - if slice is not empty, send slice down results channel
    - if tasklist is empty and treebuilding channel is closed, end search WG

5. start GR to wait for completion of search WG. when all search loops have completed, close channel and decrement master WG

6. start GR for displaying results. when search channel is closed and results channel is empty, decrement master WG

7. back in the main thread, we wait for the master WG to wrap up, then complete.