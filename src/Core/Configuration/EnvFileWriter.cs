using System;
using System.Text;

namespace Modm.Configuration
{
	public class EnvFileWriter
	{
        public IDictionary<string, string> Items { get; } = new Dictionary<string, string>();


        public EnvFileWriter(IReadOnlyDictionary<string, string> items)
        {
            Items = new Dictionary<string, string>(items);
        }

        public void Add(string key, string value)
        {
            if (!Items.TryAdd(key, value))
            {
                Items[key] = value;
            }
        }

        /// <summary>
        /// Writes all items to the specified .env file
        /// </summary>
        /// <param name="path">the env file path</param>
        /// <returns></returns>
        public async Task WriteAsync(string path)
        {
            if (Items.Count == 0)
            {
                return;
            }

            await File.AppendAllLinesAsync(path, Items.Select(pair => $"{pair.Key}={pair.Value}"), Encoding.UTF8);
        }
    }
}

