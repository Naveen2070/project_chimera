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
