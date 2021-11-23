# DISYS_Mini_Project_3

## How to run
### Run servers
- Open 1 to 90 terminal windows in DISYS_Mini_Project_3/Server/ and run the following command in each to run the servers: 
```console
go run .
```
- You should open atleast two to try and crash one of them by closing the terminal window or clicking Crtl + C on Windows.
### Run clients
- Open as many terminal windows in DISYS_Mini_Project_3/Client/ as you would like and type the following command in each to run the clients:
```console
go run .
```
## How to auction
- To make a bid, simply type an integer and press enter.
- To query the system in order to know the state of the auction, enter "r" and press enter.

## Good to know
- The action will start when the first bid is made.
- The action will last 1 minut.
- Make sure to start all the servers first as the clients will ping the servers and find all the available servers.

