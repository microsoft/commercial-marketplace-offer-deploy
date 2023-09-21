//using System.Net.Http;
//using Microsoft.Extensions.Options;
//using Modm.Engine;
//using Modm.Engine.Jenkins.Client;
//using Moq;

//namespace Modm.Tests.Core.Engine.Jenkins.Client
//{
//    public class JenkinsClientFactoryTests
//    {
//        [Fact]
//        public async Task Create_JenkinsClientWithRetryPolicy()
//        {
//            // Arrange
//            var mockClientFactory = new Mock<IHttpClientFactory>();
//            var mockApiTokenClient = new Mock<ApiTokenClient>();
//            var mockOptions = new Mock<IOptions<JenkinsOptions>>();
//            mockOptions.Setup(x => x.Value).Returns(new JenkinsOptions
//            {
//                BaseUrl = "http://jenkins.example.com",
//                UserName = "testuser",
//                Password = "test",
//                ApiToken = "test"
//            });

//            var httpClient = new HttpClient();
//            mockClientFactory.Setup(factory => factory.CreateClient(It.IsAny<string>())).Returns(httpClient);

//            var jenkinsClientFactory = new JenkinsClientFactory(mockClientFactory.Object, mockApiTokenClient.Object, mockOptions.Object);

//            // Act
//            var jenkinsClient = await jenkinsClientFactory.Create();

//            // Assert
//            Assert.NotNull(jenkinsClient);

//            // Verify that the HttpClient has a retry policy attached
//            var response = await httpClient.GetAsync("http://jenkins.example.com");
//            Assert.True(response.IsSuccessStatusCode);
//        }
//    }
//}

