# dogmatic

An experimental utility for interacting with DGT's Electronic Chess
Boards, written in Go. The development platform is Linux; on other
platforms, YMMV.

This version works much better than
[godgt](http://github.com/kgigitdev/godgt). That version attempted to
parse individual low-level events and work out what moves were
intended using a complicated heuristic of keeping track of what pieces
were "in flight".

This version does away with all of that.

It now uses the excellent
[notnil/chess](http://github.com/notnil/chess) library to generate all
legal next positions by brute force (which is, in reality, not a great
number of moves), and compares each new position with all legal next
positions to see which to accept.

## Installation

```
go get github.com/kgigitdev/dogmatic
cd ${GOPATH}/src/github.com/kgigitdev/dogmatic/dgtcli
go install
```

## Dependencies

```
go get github.com/jacobsa/go-serial/serial
go get github.com/jessevdk/go-flags
go get github.com/notnil/chess
```

## Information

Luckily, the DGT boards don't need any additional driver under Linux:
plugging one in should automatically load the FTDI USB-serial driver
(ftdi_sio and usbserial) on a reasonably up-to-date Linux
installation. You should also see a serial device like one of the
following appear:

```
$ ls -la /dev/ttyUSB*
crw-rw---- 1 root dialout 188, 0 Apr 11 18:52 /dev/ttyUSB0

$ ls -la /dev/ttyACM*
crw-rw---- 1 root dialout 188, 0 Apr 11 18:52 /dev/ttyACM0
```

In order to make this readable by you, you might need to add yourself
to the group owning the device, using a command like the following,
after which you will need to log out and back in.

```
sudo usermod -a -G dialout ${USER}
```

## Example usage

```
$ dgtcli -d /dev/ttyUSB0 
2018/08/12 13:00:47 Opening port ...
2018/08/12 13:00:47 Creating board ...
2018/08/12 13:00:47 Starting byte reader ...
2018/08/12 13:00:47 Starting command processor ...
2018/08/12 13:00:47 Resetting board ...
2018/08/12 13:00:47 New Game Started!
2018/08/12 13:00:54 Not a move!
2018/08/12 13:00:54 In e2, want ♙ have  
2018/08/12 13:00:54 Accepted move: e2e4
rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1

1.e4 *

 A B C D E F G H
8♜ ♞ ♝ ♛ ♚ ♝ ♞ ♜ 
7♟ ♟ ♟ ♟ ♟ ♟ ♟ ♟ 
6- - - - - - - - 
5- - - - - - - - 
4- - - - ♙ - - - 
3- - - - - - - - 
2♙ ♙ ♙ ♙ - ♙ ♙ ♙ 
1♖ ♘ ♗ ♕ ♔ ♗ ♘ ♖ 

2018/08/12 13:01:01 Not a move!
2018/08/12 13:01:01 In e7, want ♟ have  
2018/08/12 13:01:02 Accepted move: e7e5
rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w KQkq e6 0 2

1.e4 e5  *

 A B C D E F G H
8♜ ♞ ♝ ♛ ♚ ♝ ♞ ♜ 
7♟ ♟ ♟ ♟ - ♟ ♟ ♟ 
6- - - - - - - - 
5- - - - ♟ - - - 
4- - - - ♙ - - - 
3- - - - - - - - 
2♙ ♙ ♙ ♙ - ♙ ♙ ♙ 
1♖ ♘ ♗ ♕ ♔ ♗ ♘ ♖ 
```

## Bugs and Known Issues

* Sliding piece moves can confuse the program. For example, when
  executing a "Bb5" move, a player might slide the bishop over d2, e3
  and c4, and if the board happens to scan the piece in that position,
  then Bd3/Be4/Bd4 becomes the accepted move, and the later "Bb5" is
  rejected as being an invalid move (because the same player is
  seemingly playing twice).
* We generate legal moves on each and every piece event, which is
  inefficient.
* The output is too verbose.

