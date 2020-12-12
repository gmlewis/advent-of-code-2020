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
  int de = 0;
  int dn = 0;
  int epos = 0;
  int npos = 0;
  List<String> lines = [];

  Puzzle(String filename) {
    lines = File(filename).readAsLinesSync();
    de = 1;
  }

  iterate() {
    for (var line in lines) {
      var amt = int.parse(line.substring(1));
      switch (line.substring(0, 1)) {
        case 'N':
          npos += amt;
          break;
        case 'S':
          npos -= amt;
          break;
        case 'E':
          epos += amt;
          break;
        case 'W':
          epos -= amt;
          break;
        case 'L':
          while (amt > 0) {
            var tmp = dn;
            dn = de;
            de = -tmp;
            amt -= 90;
          }
          break;
        case 'R':
          while (amt > 0) {
            var tmp = dn;
            dn = -de;
            de = tmp;
            amt -= 90;
          }
          break;
        case 'F':
          npos += amt * dn;
          epos += amt * de;
          break;
      }
      // print('$line: ($epos,$npos) ($de,$dn)');
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
