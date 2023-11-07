using System;
namespace ClientApp.Security
{
	/// <summary>
	/// We will only have a single, admin user
	/// </summary>
	public record AuthenticatedUser
	{
		/// <summary>
		/// The id that represents the current, authenticated user
		/// </summary>
		public Guid Id { get; set; }
		public DateTimeOffset Expires { get; set; }
		public string Token { get; set; }

        public AuthenticatedUser()
		{
			Id = Guid.NewGuid();
			Expires = DateTimeOffset.UtcNow.AddMinutes(30);
		}
	}
}

