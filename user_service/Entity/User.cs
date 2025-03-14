namespace user_service.Entity
{
    using System;
    using System.ComponentModel.DataAnnotations;
    using System.ComponentModel.DataAnnotations.Schema;
    using Microsoft.EntityFrameworkCore;

    namespace AuthService.Model
    {
        [Table("users")]
        [Index(nameof(Email), IsUnique = true)]
        [Index(nameof(Username), IsUnique = true)]
        public class User: AuditableEntity
        {
            [Key]
            [DatabaseGenerated(DatabaseGeneratedOption.Identity)]
            [Column("id")]
            public Guid Id { get; set; }

            [Required]
            [Column("first_name")]
            public string FirstName { get; set; } = null!;

            [Required]
            [Column("last_name")]
            public string LastName { get; set; } = null!;

            [Required]
            [Column("email")]
            public string Email { get; set; } = null!;

            [Required]
            [Column("username")]
            public string Username { get; set; } = null!;

            [Column("password")]
            public string Password { get; set; } = null!;

            [Column("role")]
            public string Role { get; set; } = null!;

            [Required]
            [Column("status")]
            public int Status { get; set; } = 1;

            [Column("created_on")]
            public override DateTime CreatedOn { get; set; } = DateTime.UtcNow;

            [Column("updated_on")]
            public override DateTime UpdatedOn { get; set; } = DateTime.UtcNow;
        }
    }
}
