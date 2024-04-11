
using System.Text.Json;
using static System.Text.Json.JsonElement;

namespace Modm.Tests.Resources
{
    /// <summary>
    /// The resource representation for the templates/mainTemplate.json file
    /// so we can assert certain conditions against its contents
    /// </summary>
    public class MainTemplate
    {
        public readonly string Content;
        public readonly JsonDocument Json;

        private MainTemplate()
        {
            var path = Path.Combine(AppDomain.CurrentDomain.BaseDirectory, "../../../../", "templates/mainTemplate.json");
            this.Content = File.ReadAllText(path);
            this.Json = JsonDocument.Parse(Content);
        }

        public static MainTemplate Get()
        {
            return new MainTemplate();
        }

        public List<JsonElement> GetResourcesByCommonTag()
        {
            var result = new List<JsonElement>();
            var filter = (JsonElement e) =>
            {
                return e.ValueKind == JsonValueKind.Object &&
                e.EnumerateObject().Any(p =>
                    p.Name == "tags" &&
                    p.Value.ValueKind == JsonValueKind.String &&
                    p.Value.GetString() == "[variables('commonTags')]");
            };

            Find(Json.RootElement.GetProperty("resources").EnumerateArray(), filter, result);

            return result;
        }

        private static void Find(ArrayEnumerator array, Func<JsonElement, bool> predicate, List<JsonElement> result)
        {
            foreach (var element in array)
            {
                if (element.ValueKind == JsonValueKind.Array)
                {
                    Find(element.EnumerateArray(), predicate, result);
                    continue;
                }

                if (element.ValueKind == JsonValueKind.Object)
                {
                    if (predicate(element))
                    {
                        result.Add(element);
                    }
                }
            }
        }
    }
}

