using System.Text.Json;
using Modm.Deployments;
using Xunit;

public class ArmParameterExtensionsTests
{
    [Fact]
    public void ArmParametersJson_Should_HaveExpectedValues()
    {
        // Arrange
        var inputParameters = new Dictionary<string, object>
        {
            {"siteName", "GEN-UNIQUE"},
            {"administratorLogin", "GEN-UNIQUE-12"},
            {"administratorLoginPassword", "GEN-PASSWORD"}
        };

        var expectedJson = @"
        {
            ""parameters"": {
                ""siteName"": { ""value"": ""GEN-UNIQUE"" },
                ""administratorLogin"": { ""value"": ""GEN-UNIQUE-12"" },
                ""administratorLoginPassword"": { ""value"": ""GEN-PASSWORD"" }
            }
        }";

        // Act
        var actualJson = inputParameters.ToArmParametersJson();
        var actualDoc = JsonDocument.Parse(actualJson);
        var expectedDoc = JsonDocument.Parse(expectedJson);

        // Assert
        Assert.Equal(
            expectedDoc.RootElement.GetProperty("parameters").GetProperty("siteName").GetProperty("value").GetString(),
            actualDoc.RootElement.GetProperty("parameters").GetProperty("siteName").GetProperty("value").GetString()
        );

        Assert.Equal(
            expectedDoc.RootElement.GetProperty("parameters").GetProperty("administratorLogin").GetProperty("value").GetString(),
            actualDoc.RootElement.GetProperty("parameters").GetProperty("administratorLogin").GetProperty("value").GetString()
        );

        Assert.Equal(
            expectedDoc.RootElement.GetProperty("parameters").GetProperty("administratorLoginPassword").GetProperty("value").GetString(),
            actualDoc.RootElement.GetProperty("parameters").GetProperty("administratorLoginPassword").GetProperty("value").GetString()
        );
    }
}
