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
  cpu.execute();
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

  execute() {
    var seen = Set<int>();

    while (true) {
      if (seen.contains(ip)) {
        print('Solution Accumulator = $accumulator');
        return;
      }
      seen.add(ip);

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
