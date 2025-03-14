namespace user_service.Utils.Enums
{
    public static class UserRole
    {
        public const string Admin = "ADMIN";
        public const string User = "USER";
        public const string Guest = "GUEST";
        public const string Super = "SUPER";

        // Optional: You can also add a method to get all roles as a list
        public static IEnumerable<string> GetAllRoles()
        {
            return new List<string> { Admin, User, Guest, Super };
        }
    }
}
