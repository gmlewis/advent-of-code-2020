// -*- compile-command: "dart main.dart ../example1.txt ../input.txt"; -*-

import 'dart:io';

main(List<String> args) {
  for (var arg in args) {
    process(arg);
  }

  print('Done.');
}

class Instruction {
  String op;
  int arg;
  Instruction({this.op, this.arg});

  toString() => 'op: $op, arg: $arg';
}

process(String filename) {
  print('Processing $filename ...');

  var program = List<Instruction>();
  var lines = File(filename).readAsLinesSync();
  for (var line in lines) {
    var parts = line.split(' ');
    var ins = Instruction(op: parts[0], arg: int.parse(parts[1]));
    program.add(ins);
  }

  var cpu = CPU(program: program);
  var recording = cpu.execute()[0];

  for (var i = recording.length-1; i >= 0; i--) {
    var shallowCopy = List<Instruction>.from(program);
    cpu = CPU(program: shallowCopy);
    var ip = recording[i];
    var ins = program[ip];
    if (ins.op == 'nop') {
      print('Experiment: Changing op at $ip from nop to jmp');
      cpu.program[ip].op = 'jmp';
    } else if (ins.op == 'jmp') {
      print('Experiment: Changing op at $ip from jmp to nop');
      cpu.program[ip].op = 'nop';
    } else {
      continue;
    }

    if (cpu.execute()[1]) {
      break;
    }
  }
}

class CPU {
  int accumulator;
  int ip;
  List<Instruction> program;
  CPU({
      this.accumulator = 0,
      this.ip = 0,
      this.program,
  });

  List execute() {
    var seen = Set<int>();
    var recording = List<int>();

    while (true) {
      if (seen.contains(ip)) {
        print('Infinite loop: Accumulator = $accumulator');
        return [recording, false];
      }
      if (ip >= program.length) {
        print('Normal termination: Accumulator = $accumulator');
        return [recording, true];
      }
      seen.add(ip);

      recording.add(ip);
      var ins = program[ip];
      switch (ins.op) {
        case 'nop':
        ip++;
        break;
        case 'acc':
        accumulator += ins.arg;
        ip++;
        break;
        case 'jmp':
        ip += ins.arg;
        break;
      }
    }
  }
}
