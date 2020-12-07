// -*- compile-command: "dart main.dart ../example1.txt ../input.txt"; -*-

import 'dart:io';

main(List<String> args) {
  for (var arg in args) {
    process(arg);
  }

  print('Done.');
}

class Contains {
  int quant;
  String color;
  Contains({this.quant, this.color});

  toString() => 'Contains(quant: $quant, color: $color)';
}

process(String filename) {
  print('Processing $filename ...');

  var lineRE = RegExp(r'^(.*?) bags contain (.*)\.$');
  var containsRE = RegExp(r',?\s*(\d+) (.*?) bags?');

  var rules = Map<String,List<Contains>>();

  var lines = File(filename).readAsLinesSync();
  for (var line in lines) {
    var m = lineRE.firstMatch(line);
    // print('${m.group(0)}; ${m.group(1)}; ${m.group(2)}');
    if (m.group(2) == 'no other bags') {
      rules[m.group(1)] = [];  // empty list
      // print('No bags: rules[${m.group(1)}] = []');
      continue;
    }

    var m2 = containsRE.allMatches(m.group(2));
    m2.forEach((f) {
        var q = int.parse(f.group(1));
        // print('${f.group(0)}; ${f.group(1)}; ${f.group(2)}');
        var c = Contains(quant: q, color: f.group(2));
        // if (rules.containsKey(m.group(1))) {
        //   print('BEFORE: rules[${m.group(1)}] = ${rules[m.group(1)]}');
        // } else {
        //   print('BEFORE: rules[${m.group(1)}] = null');
        // }
        rules.update(m.group(1), (cs) {
            // print('CS BEFORE: cs=$cs');
            cs.add(c);
            // print('CS AFTER: cs=$cs');
            return cs;  // CRITICAL TO PROVIDE!!!
          }, ifAbsent: () => [c]);
        // print('AFTER: rules[${m.group(1)}] = ${rules[m.group(1)]}');
    });
  }

  var count = 0;
  rules.keys.forEach((color) {
      if (canContain(color, "shiny gold", rules)) count++;
  });

  print('Solution: $count');
}

bool canContain(String key, String test, Map<String,List<Contains>> rules, {Map<String,bool> seen}) {
  if (seen == null) {
    seen = {};
  }

  if (!rules.containsKey(key) || rules[key] == null) return false;

  return rules[key].any((v) {
      if (v.color == test) return true;

      if (seen.containsKey(v.color)) return false;
      seen[v.color] = true;

      return canContain(v.color, test, rules, seen: seen);
  });
}
