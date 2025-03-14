namespace user_service.Utils.Enums
{
    public static class UserStatus
    {
        public const int Active = 1;
        public const int Inactive = 0;
        public const int Locked = 2;
        public const int Expired = 3;

        // Optional: You can also add a method to get all statuses as a dictionary
        public static IDictionary<int, string> GetAllStatuses()
        {
            return new Dictionary<int, string>
            {
                { Active, "Active" },
                { Inactive, "Inactive" },
                { Locked, "Locked" },
                { Expired, "Expired" }
            };
        }
    }
}
