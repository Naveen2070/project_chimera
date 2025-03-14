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
