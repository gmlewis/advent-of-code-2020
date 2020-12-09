// -*- compile-command: "dart main.dart -p 5 ../example1.txt && dart main.dart -p 25 ../input.txt"; -*-

import 'dart:io';

main(List<String> args) {
  if (args.length != 3 || args[0] != '-p') {
    print('usage: dart main.dart -p 25 ../imput.txt');
    return;
  }

  var preamble = int.parse(args[1]);
  process(args[2], preamble);

  print('Done.');
}

process(String filename, int preamble) {
  print('Processing $filename ...');

  var ints = List<int>();
  var lines = File(filename).readAsLinesSync();
  for (var line in lines) {
    var n = int.parse(line);
    ints.add(n);
  }

  for (var i = preamble; i < ints.length; i++) {
    if (foundSolution(i, ints, preamble)) {
      break;
    }
  }
}

bool foundSolution(int n, List<int> ints, int preamble) {
  var v = ints[n];
  var seen = Set<int>();
  for (int i = n - preamble; i < n; i++) {
    var d = v - ints[i];
    if (seen.contains(d)) {
      return false;
    }
    seen.add(ints[i]);
  }
  print('Solution: $v');
  return true;
}
