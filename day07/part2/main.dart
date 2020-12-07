// -*- compile-command: "dart main.dart ../example1.txt ../example2.txt ../input.txt"; -*-

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

  var count = contains("shiny gold", rules);

  print('Solution: ${count-1}');
}

int contains(String key, Map<String,List<Contains>> rules) {
  var total = 1;

  if (!rules.containsKey(key) || rules[key] == null) return 1;

  rules[key].forEach((v) {
      var count = contains(v.color, rules);
      total += (count * v.quant);
  });

  return total;
}
