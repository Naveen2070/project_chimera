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
using user_service.Utils.Enums;

namespace user_service.Utils.Validators
{
    public class ValidUserRoleAttribute : ValidationAttribute
    {
        protected override ValidationResult? IsValid(object? value, ValidationContext validationContext)
        {
            if (value == null)
                return new ValidationResult("Role is required.");

            var roleValue = value.ToString();

            if (!UserRole.GetAllRoles().Contains(roleValue))
            {
                return new ValidationResult($"Role must be one of: {string.Join(", ", UserRole.GetAllRoles())}");
            }

            return ValidationResult.Success;
        }
    }
}
