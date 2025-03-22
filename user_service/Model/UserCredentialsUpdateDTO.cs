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
using System.ComponentModel.DataAnnotations;

namespace user_service.Model
{
    public class UserCredentialsUpdateDTO
    {
        [Required]
        public Guid UserId { get; set; }

        [Required]
        [StringLength(100, MinimumLength = 6, ErrorMessage = "Old password must be at least 6 characters long.")]
        public string OldPassword { get; set; } = null!;

        [Required]
        [StringLength(100, MinimumLength = 6, ErrorMessage = "New password must be at least 6 characters long.")]
        [RegularExpression(@"^(?=.*[a-z])(?=.*[A-Z])(?=.*\d).+$",
            ErrorMessage = "New password must contain at least one uppercase letter, one lowercase letter, and one number.")]
        public string NewPassword { get; set; } = null!;

        [Required]
        [Compare("NewPassword", ErrorMessage = "Confirm password doesn't match.")]
        public string ConfirmPassword { get; set; } = null!;
    }
}
