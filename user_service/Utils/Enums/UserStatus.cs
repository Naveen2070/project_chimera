//	Copyright 2025 Naveen R
//
//		Licensed under the Apache License, Version 2.0 (the "License");
//		you may not use this file except in compliance with the License.
//		You may obtain a copy of the License at
//
//		http://www.apache.org/licenses/LICENSE-2.0
//
//		Unless required by applicable law or agreed to in writing, software
//		distributed under the License is distributed on an "AS IS" BASIS,
//		WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//		See the License for the specific language governing permissions and
//		limitations under the License.
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
