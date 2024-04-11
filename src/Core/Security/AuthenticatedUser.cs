using System;

namespace Modm.Security
{
	/// <summary>
	/// We will only have a single, admin user. represents a combined
	/// instance of both an authenticated user and a unique session, keyed of the Id
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

