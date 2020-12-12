// -*- compile-command: "dart main.dart ../example1.txt ../input.txt"; -*-

import 'dart:io';

main(List<String> args) {
  for (var arg in args) {
    process(arg);
  }

  print('Done.');
}

process(String filename) {
  print('Processing $filename ...');

  var puz = Puzzle(filename);

  puz.iterate();

  var dist = puz.manhattan();

  print('Solution: $dist');
}

class Puzzle {
  int wepos = 0;
  int wnpos = 0;
  int epos = 0;
  int npos = 0;
  List<String> lines = [];

  Puzzle(String filename) {
    lines = File(filename).readAsLinesSync();
    wepos = 10;
    wnpos = 1;
  }

  iterate() {
    for (var line in lines) {
      var amt = int.parse(line.substring(1));
      switch (line.substring(0, 1)) {
        case 'N':
          wnpos += amt;
          break;
        case 'S':
          wnpos -= amt;
          break;
        case 'E':
          wepos += amt;
          break;
        case 'W':
          wepos -= amt;
          break;
        case 'L':
          while (amt > 0) {
            var tmp = wnpos;
            wnpos = wepos;
            wepos = -tmp;
            amt -= 90;
          }
          break;
        case 'R':
          while (amt > 0) {
            var tmp = wnpos;
            wnpos = -wepos;
            wepos = tmp;
            amt -= 90;
          }
          break;
        case 'F':
          npos += amt * wnpos;
          epos += amt * wepos;
          break;
      }
      // print('$line: ($epos,$npos) ($wepos,$wnpos)');
    }
  }

  int manhattan() {
    var x = epos;
    if (x < 0) {
      x = -x;
    }
    var y = npos;
    if (y < 0) {
      y = -y;
    }
    return x + y;
  }
}
