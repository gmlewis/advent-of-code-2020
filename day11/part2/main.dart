// -*- compile-command: "pub get && dart main.dart ../example1.txt ../input.txt"; -*-

import 'dart:io';

import 'package:collection/equality.dart';
import 'package:tuple/tuple.dart';

main(List<String> args) {
  for (var arg in args) {
    process(arg);
  }

  print('Done.');
}

process(String filename) {
  print('Processing $filename ...');

  var puz = Puzzle(filename);

  while (true) {
    var newGrid = puz.iterate();
    if (MapEquality().equals(puz.grid, newGrid)) {
      break;
    }
    puz.grid = newGrid;
  }

  var occupied = puz.occupied();

  print('Solution: $occupied');
}

class Puzzle {
  int width = 0;
  int height = 0;
  var grid = Map<Tuple2<int,int>,String>();

  Puzzle(String filename) {
    var lines = File(filename).readAsLinesSync();
    for (var line in lines) {
      this.width = line.length;
      var chars = line.split('');
      for (var x = 0; x < chars.length; x++) {
        if (chars[x] == 'L') {
          var key = Tuple2<int,int>(x, this.height);
          this.grid[key] = 'L';
        }
      }
      this.height++;
    }
  }

  Map<Tuple2<int,int>,String> iterate() {
    var r = Map<Tuple2<int,int>,String>();
    grid.forEach((k, v) {
        var adj = countAdjacentOccupied(k);
        if (v == 'L' && adj == 0) {
          r[k] = '#';
        } else if (v == '#' && adj >= 5) {
          r[k] = 'L';
        } else {
          r[k] = v;
        }
    });
    return r;
  }

  bool adjacent(int x, int dx, int y, int dy) {
    var u = x + dx;
    var v = y + dy;
    var adj = 0;
    while (u >= 0 && v >= 0 && u < width && v < height) {
      var key = Tuple2<int,int>(u, v);
      var val = grid[key];
      if (val == '#') {
        return true;
      }
      if (val == 'L') {
        return false;
      }
      u += dx;
      v += dy;
    }
    return false;
  }

  int countAdjacentOccupied(Tuple2<int,int> k) {
    var x = k.item1;
    var y = k.item2;
    var adj = 0;

    if (adjacent(x,-1, y,-1)) adj++;
    if (adjacent(x, 0, y,-1)) adj++;
    if (adjacent(x, 1, y,-1)) adj++;

    if (adjacent(x,-1, y, 0)) adj++;
    if (adjacent(x, 1, y, 0)) adj++;

    if (adjacent(x,-1, y, 1)) adj++;
    if (adjacent(x, 0, y, 1)) adj++;
    if (adjacent(x, 1, y, 1)) adj++;

    return adj;
  }

  int occupied() {
    var count = 0;
    grid.forEach((k, v) {
        if (v == '#') count++;
    });
    return count;
  }
}
