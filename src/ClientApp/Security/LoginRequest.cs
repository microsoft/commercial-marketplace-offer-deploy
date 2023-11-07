using System;
namespace ClientApp.Security
{
	public record LoginRequest
	{
		public string Username { get; set; }
		public string Password { get; set; }
	}
}