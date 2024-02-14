// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿using System;
using System.Linq;
using System.Text;

namespace Modm.Configuration
{
	public class EnvFile : IDisposable
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

        private EnvFile(string filePath)
        {
            this.filePath = Path.GetFullPath(filePath);
            this.items = new Dictionary<string, string>();
        }

        public static EnvFile New(string filePath)
        {
            return new EnvFile(filePath);
        }

        public async Task<bool> AnyAsync()
        {
            if (!read)
            {
               await ReadAsync();
            }
            return Items.Count > 0;
        }

        public string[] ToArray()
        {
            return Items.Select(item => $"{item.Key}={item.Value}").ToArray();
        }

        /// <summary>
        /// Set a value by key
        /// </summary>
        /// <param name="key"></param>
        /// <param name="value"></param>
        public void Set(string key, string value)
        {
            if (!items.TryAdd(key, value))
            {
                items[key] = value;
            }
        }

        /// <summary>
        /// Set a value by key
        /// </summary>
        /// <param name="key"></param>
        /// <param name="value"></param>
        public void Remove(string key)
        {
            items.Remove(key);
        }

        /// <summary>
        /// Will read the .env file and load the key value items. This is performed on the
        /// first call to read. Consecutive calls will do nothing.
        /// </summary>
        /// <returns></returns>
        public async Task ReadAsync()
        {
            if (read)
            {
                return;
            }

            if (!File.Exists(filePath))
            {
                read = true;
                return;
            }

            foreach (var line in await File.ReadAllLinesAsync(filePath))
            {
                var parts = line.Split('=', 2, StringSplitOptions.RemoveEmptyEntries);
                items.TryAdd(parts[0], parts[1]);
            }

            read = true;
        }

        /// <summary>
        /// Writes all items to the specified .env file
        /// </summary>
        /// <param name="path">the env file path</param>
        /// <returns></returns>
        public async Task SaveAsync()
        {
            if (Items.Count == 0)
            {
                return;
            }

            var lines = Items.Select(pair => $"{pair.Key}={pair.Value}");
            await File.WriteAllLinesAsync(filePath, lines, Encoding.UTF8);
        }

        public void Dispose()
        {
            items.Clear();
        }
    }
}

