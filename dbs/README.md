# What is this

I wanted to come up with a way of writing data access code so that each individual query could be
invoked by both the `*sql.DB` and `*sql.Tx` structs. This is the design that works for me.

The nice thing is that `*dbs.DB` and `*dbs.Tx` both have all queries available as receiver methods, but
only `*dbs.DB` can create a new transaction, and only `*dbs.Tx` can commit or rollback.

Additional feature of this design is that queries are lazily prepared.