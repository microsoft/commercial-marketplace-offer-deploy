using System;
using System.IO.Compression;

namespace Modm.Compression
{
    public class DirectoryZipper
    {
        public void ZipDirectory(string sourceDirectoryName, string destinationZipFileName)
        {
            var destinationDirectory = Path.GetDirectoryName(destinationZipFileName);

            if (!Directory.Exists(destinationDirectory))
            {
                Directory.CreateDirectory(destinationDirectory);
            }

            ZipFile.CreateFromDirectory(sourceDirectoryName, destinationZipFileName);
        }
    }
}

