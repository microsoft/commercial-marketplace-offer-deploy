// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿namespace Modm.Tests.Utils
{
    /// <summary>
    /// so we can test and make sure local IO is cleanedup
    /// </summary>
    /// <typeparam name="TTest"></typeparam>
    public class DisposableDirectory<TTest> : IDisposable
	{
        readonly DirectoryInfo directory;

        public string FullName => directory.FullName;
        public DirectoryInfo Info => directory;

        public DisposableDirectory()
        {
            var directoryPath = Path.Combine(Path.GetTempPath(), typeof(TTest).Name, Guid.NewGuid().ToString());
            directory = Directory.CreateDirectory(directoryPath);
        }

        public void Delete()
        {
            if (directory == null)
            {
                return;
            }

            try
            {
                directory.Delete(recursive: true);
            }
            catch
            {
            }
        }

        public void Dispose()
        {
            Delete();
        }
    }

    partial class Test
    {
        public static DisposableDirectory<TTest> Directory<TTest>()
        {
            return new DisposableDirectory<TTest>();
        }
    }
}

