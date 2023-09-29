using System;
namespace Modm.Tests.Utils
{
	partial class Test
	{
		public static class DataFile
		{
			const string DataFolderName = "Data";

			public static FileInfo Get(string fileName)
			{
                var fullPath = Path.Combine(AppDomain.CurrentDomain.BaseDirectory, DataFolderName, fileName);
				return new FileInfo(fullPath);
            }
		}
	}
}

