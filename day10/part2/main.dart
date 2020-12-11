// -*- compile-command: "dart main.dart ../example1.txt ../example2.txt ../input.txt"; -*-

import 'dart:io';

main(List<String> args) {
  for (var arg in args) {
    process(arg);
  }

  print('Done.');
}

process(String filename) {
  print('Processing $filename ...');

  var ints = Set<int>();
  var max = 0;
  var lines = File(filename).readAsLinesSync();
  for (var line in lines) {
    var n = int.parse(line);
    ints.add(n);
    if (n > max) max = n;
  }

  var cache = Map<int,int>();
  var count = countPossibilities(0, ints, max, cache);

  print('Solution: $count');
}

int countPossibilities(int lastN, Set<int> ints, int max, Map<int,int> cache) {
  if (cache.containsKey(lastN)) return cache[lastN];

  if (lastN == max) {
    cache[lastN] = 1;
    return 1;
  }

  if (lastN > max) {
    cache[lastN] = 0;
    return 0;
  }

  var result = 0;
  for (var i = 1; i <= 3; i++) {
    if (ints.contains(lastN+i)) {
      result += countPossibilities(lastN+i, ints, max, cache);
    }
  }

  cache[lastN] = result;
  return result;
}
