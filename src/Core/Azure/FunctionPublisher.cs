using System;
using System.Diagnostics;


namespace Modm.Azure
{
    public class FunctionPublisher
    {
        private readonly string workingDirectory;

        public FunctionPublisher(string workingDirectory)
        {
            this.workingDirectory = workingDirectory ?? throw new ArgumentNullException(nameof(workingDirectory));
        }

        public async Task PublishAsync(string clientAppName)
        {
            var startInfo = new ProcessStartInfo
            {
                FileName = "func",
                Arguments = $"azure functionapp publish {clientAppName}",
                RedirectStandardOutput = true,
                RedirectStandardError = true,
                UseShellExecute = false,
                CreateNoWindow = true,
                WorkingDirectory = this.workingDirectory 
            };

            using var process = new Process { StartInfo = startInfo };

            process.OutputDataReceived += (sender, args) => Console.WriteLine(args.Data);
            process.ErrorDataReceived += (sender, args) => Console.WriteLine($"Error: {args.Data}");

            process.Start();

            process.BeginOutputReadLine();
            process.BeginErrorReadLine();

            await process.WaitForExitAsync();
        }
    }
}

