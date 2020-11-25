## Server-Commands
* ```host (string)```: 
    flag for the hostname of the server
* ```port (int)```: 
    flag for the port of the server
* ```servername (string)```: 
    flag for the name of server
    
## Worker-Commands
* ```local bind hostname (string)```: 
        bind address to (default "localhost")
* ```local bind port (int)```: 
        bind address to (default 8081")      
* ```remote hostname (string)```:
        remote address (default "localhost")
* ```remote port (int)```:
        remote address (default 8080")
* ```createtree (bool)```:
        create tree. default not creating tree
* ```deletetree (bool)```:
        flag for deleting a tree
* ```traversetree (bool)```:
        flag for traversing the tree

* ```deletekey (bool)```: 
        flag for deleting a key/value in the tree 
* ```findvalue (bool)```:
        flag for finding a value in the tree
* ```insertvalue (bool)```:
        flag for inserting a value to the tree
* ```leafsize (int)```:
        leafsize value when create a tree

* ```id (int)```:
        flag for id. necessary for all operations
* ```token (string)```:
        flag for token. necessary for all operations
* ```key (int)```:
        key when inserting/deleting/finding values
* ```value (string)```:
        value when inserting a key

## Usage
- First, start the tree-server:

    ```
        cd tree-server && go run main.go
    ```

- Second, start the tree-worker(client):

    ```
        cd tree-worker && go run main.go
    ```
Without any of the available flags, it will print out the currently available trees!

- Show all available trees: 
    ```
    go run main.go
    ```
- Create a tree: 
    ```
    go run main.go --createtree --leafsize=100
    ```
- Delete a tree: 
    ```
    go run main.go--deletetree --token="xxxxxx" --id=19937
    ```
- Traverse a tree:  
    ```
    go run main.go --traversetree --token="xxxxxx" --id=19937
    ```
- Insert key-value to a tree: 
    ```
    go run main.go --insert --token="xxxxxx" --id=19937 --key=123 --value="google"
    ```
- Find value by key in a tree: 
    ```
    go run main.go --find --token="xxxxxx" --id=19937 --key=5
    ```
- Delete key in a tree: 
    ```
    go run main.go --deletekey --token="xxxxxx" --id=19937 --key=123
    ```
