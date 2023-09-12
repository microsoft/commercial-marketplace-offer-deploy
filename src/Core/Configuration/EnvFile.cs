using System;
namespace Modm.Configuration
{
    public class EnvFile
    {
        private readonly string filePath;
        readonly Dictionary<string, string> items;
        bool read;

        public IReadOnlyDictionary<string, string> Items
        {
            get
            {
                return items;
            }
        }

        public bool HasItems
        {
            get
            {
                if (!read)
                {
                    Read();
                }
                return items.Count > 0;
            }
        }

        private EnvFile(string filePath)
        {
            this.items = new Dictionary<string, string>();
            this.filePath = filePath;

            if (!File.Exists(filePath))
            {
                throw new InvalidOperationException("File path is invalid.");
            }
        }

        public static EnvFile FromPath(string filePath)
        {
            return new EnvFile(filePath);
        }

        public void Read()
        {
            foreach (var line in File.ReadAllLines(filePath))
            {
                var parts = line.Split('=', 2, StringSplitOptions.RemoveEmptyEntries);
                items.Add(parts[0], parts[1]);
            }

            read = true;
        }

        /// <summary>
        /// Causes all key value items read from the .env file to be applied to the environment
        /// </summary>
        public void SetEnvironmentVariables()
        {
            if (!read)
            {
                Read();
            }

            foreach (var item in Items)
            {
                Environment.SetEnvironmentVariable(item.Key, item.Value);
            }
        }
    }
}

