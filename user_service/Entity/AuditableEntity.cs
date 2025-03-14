namespace user_service.Entity
{
    public abstract class AuditableEntity
    {
        public virtual DateTime CreatedOn { get; set; }
        public virtual DateTime UpdatedOn { get; set; }
        //public string CreatedBy { get; set; } = null!;
        //public string UpdatedBy { get; set; } = null!;
    }
}
