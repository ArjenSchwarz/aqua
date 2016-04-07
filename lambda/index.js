var child_process = require('child_process');

exports.handler = function(event, context) {
  console.log(event)
  var items = event.body.split("&");

  var args = ["--json=true"]
  for (var i = items.length - 1; i >= 0; i--) {
    args.push("--" + items[i])
  }
  console.log(args)
  var proc = child_process.spawn('./aqua', args, { stdio: [process.stdin, 'pipe', 'pipe'] });

  proc.stdout.on('data', function(line){
    var msg = JSON.parse(line);
    console.log("stdout: " + msg)
    context.succeed(msg);
  })

  proc.stderr.on('data', function(line){
    var msg = new Error(line)
    console.log("stderr: " + msg)
    context.fail(msg);
  })

  proc.on('exit', function(code){
    console.error('exit: %s', code)
    context.fail("No results")
  })
}