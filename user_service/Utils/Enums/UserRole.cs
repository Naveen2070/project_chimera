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
