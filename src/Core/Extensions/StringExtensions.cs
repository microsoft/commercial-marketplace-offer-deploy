using System;
namespace Modm.Extensions
{
    public static class StringExtensions
    {
        public static string SubstringAfterMarker(this string inputString, string marker)
        {
            int markerIndex = inputString.IndexOf(marker);

            if (markerIndex == -1)
            {
                return inputString;
            }

            return inputString.Substring(markerIndex + marker.Length);
        }
    }
}

