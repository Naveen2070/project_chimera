using System.ComponentModel.DataAnnotations;
using user_service.Utils.Enums;
using user_service.Utils.Validators;

namespace user_service.Model
{
    public class UserUpdateDTO
    {
        [Required]
        [StringLength(50, MinimumLength = 2, ErrorMessage = "First Name must be between 2 and 50 characters.")]
        public string FirstName { get; set; } = null!;

        [Required]
        [StringLength(50, MinimumLength = 2, ErrorMessage = "Last Name must be between 2 and 50 characters.")]
        public string LastName { get; set; } = null!;

        [Required]
        [EmailAddress(ErrorMessage = "Invalid email format.")]
        public string Email { get; set; } = null!;

        [Required]
        [RegularExpression(@"^[a-zA-Z0-9_]{3,20}$", ErrorMessage = "Username must be alphanumeric and 3-20 characters long.")]
        public string Username { get; set; } = null!;

        [ValidUserRole]
        [Required]
        public string Role { get; set; } = UserRole.User;

        [Range(0, 1, ErrorMessage = "Status must be either 0 (inactive) or 1 (active).")]
        public int Status { get; set; } = 1;
    }
}
