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
  print('Part 1 solution: $v');

  for (int i = n - 1; i >= preamble-2; i--) {
    var sum = 0;
    var min = ints[i];
    var max = ints[i];
    for (int j = 0; j < preamble-1; j++) {
      var sample = ints[i-j];
      if (sample < min) { min = sample; }
      if (sample > max) { max = sample; }

      sum += sample;
      if (sum > v) { break; }
      if (sum == v) {
        print('Part 2 solution: $min + $max = ${min+max}');
        return true;
      }
    }
  }

  return true;
}
