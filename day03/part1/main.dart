// -*- compile-command: "dart main.dart ../example1.txt ../input.txt"; -*-

import 'dart:io';
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
  puz.count(3, 1);
}

class Puzzle {
  int width = 0;
  int height = 0;
  Set<Tuple2<int,int>> grid = Set<Tuple2<int,int>>();

  Puzzle(String filename) {
    var lines = File(filename).readAsLinesSync();
    for (var line in lines) {
      this.width = line.length;
      var chars = line.split('');
      for (var x = 0; x < chars.length; x++) {
        if (chars[x] == '#') {
          this.grid.add(Tuple2<int,int>(x, this.height));
        }
      }
      this.height++;
    }
  }

  int count(int right, int down) {
    var posX = 0, posY = 0, count = 0;
    for (var y = 0; y < this.height; y++) {
      if (this.grid.contains(Tuple2<int,int>(posX % this.width, posY))) {
        count++;
      }
      posX += right;
      posY += down;
    }
    print('Found $count trees');
    return count;
  }
}
