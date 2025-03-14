namespace user_service.Utils
{
    public class PasswordEncoder
    {
        private readonly int _workFactor;

        public PasswordEncoder(int workFactor)
        {
            _workFactor = workFactor;
        }

        public string Encode(string password)
        {
            return BCrypt.Net.BCrypt.HashPassword(password, _workFactor);
        }

        public bool Verify(string password, string hashedPassword)
        {
            return BCrypt.Net.BCrypt.Verify(password, hashedPassword);
        }
    }
}
