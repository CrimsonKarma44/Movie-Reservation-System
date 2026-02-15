# Vercel Deployment - Troubleshooting & FAQ

## 🔧 Troubleshooting Matrix

### Scenario 1: Deployment Fails During Build

**Symptoms:**
- "Build failed" message
- Red X on deployment
- Error in build logs

**Likely Causes & Solutions:**

| Cause | Solution |
|-------|----------|
| Go version mismatch | Check `go.mod` version matches Vercel's supported versions |
| Missing dependencies | Run `go mod tidy` locally, commit, push |
| Import path wrong | Verify all imports use correct module path |
| Environment variables in code | Don't reference ENV vars in build-time code |

**Debug Steps:**
1. Click "Deployments" in Vercel
2. Click failed deployment
3. Click "Logs" tab
4. Look for error message
5. Fix locally and push again

---

### Scenario 2: Deployment Succeeds But API Returns 502 Error

**Symptoms:**
- Deployment shows ✓ Ready
- API calls return 502 Bad Gateway
- Or blank error response

**Likely Causes & Solutions:**

| Cause | Solution |
|-------|----------|
| Database connection failed | Check DB environment variables |
| Missing environment variables | Verify all vars are set in Vercel → Settings |
| Database not accessible | Check database is online, firewall allows access |
| Application panic | Check logs for stack trace |
| Wrong port | Already fixed in this deployment (uses PORT env var) |

**Debug Steps:**
1. Go to Deployments → Click latest
2. Click "Logs" tab
3. Look for error messages
4. Check "Database connected" message appears
5. If missing, environment variable issue

**Fix:**
1. Verify environment variables in Settings → Environment Variables
2. Check spelling exactly matches
3. Check values are correct
4. Click deployment → "Redeploy" button

---

### Scenario 3: Database Connection Timeout

**Symptoms:**
- Deployment seems to hang
- "connection timed out" error
- "i/o timeout" in logs

**Likely Causes & Solutions:**

| Cause | Solution |
|-------|----------|
| Database host unreachable | Ensure DB host is public IP, not localhost |
| Firewall blocking Vercel | Add firewall rule to allow Vercel IPs |
| Database offline | Check database status in provider dashboard |
| Network connectivity | Try connecting to DB from your machine |
| Wrong host | Copy-paste hostname again carefully |

**Debug Steps:**
1. Test database locally: `psql postgresql://user:pass@host/db`
2. If works locally but not on Vercel, it's a network issue
3. Contact your database provider's support

---

### Scenario 4: Authentication Not Working

**Symptoms:**
- Registration returns error
- Login fails
- Token validation fails

**Likely Causes & Solutions:**

| Cause | Solution |
|-------|----------|
| JWT secrets not set | Verify JWT_SECRET_KEY_ACCESS and JWT_SECRET_KEY_REFRESH |
| JWT secrets empty | Ensure values are not blank or "undefined" |
| Secrets changed | Existing tokens won't work with new secrets |
| Wrong secret format | Should be random string, not URL or special format |

**Debug Steps:**
1. Check Settings → Environment Variables
2. Verify `JWT_SECRET_KEY_ACCESS` is set
3. Verify `JWT_SECRET_KEY_REFRESH` is set
4. Neither should be blank
5. Values should be different from each other

**Fix:**
1. Generate new secrets: `openssl rand -base64 32`
2. Update both JWT vars
3. Redeploy

---

### Scenario 5: First Request Very Slow

**Symptoms:**
- First API request takes 5-10 seconds
- Subsequent requests are fast
- Happens every time URL isn't accessed for a while

**This is Normal!** ✓

**Explanation:**
- Serverless functions have "cold starts"
- First request spins up the container (takes 5-10 sec)
- Subsequent requests reuse the running container (<1 sec)
- Normal behavior on all serverless platforms

**Not a Problem** - but if you want to minimize:
- Add warm-up requests (make dummy request every 5 minutes)
- Vercel Pro includes "sleep guard" to keep functions warm
- Consider dedicated instance if cold starts critical

---

### Scenario 6: Environment Variables Not Working

**Symptoms:**
- Application crashes with "missing required environment"
- Database error about null values
- JWT operations fail

**Likely Causes & Solutions:**

| Cause | Solution |
|-------|----------|
| Variables not set | Add them in Settings → Environment Variables |
| Set for wrong environment | Ensure "Production" is selected, not just "Preview" |
| Preview deploy not showing | Preview deployments need separate variable config |
| Typos in names | Must match exactly (case-sensitive on Linux) |
| Values with spaces | Wrap in quotes if needed |

**Debug Steps:**
1. Go to Settings → Environment Variables
2. Scroll through entire list
3. Check all required vars are present:
   - DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME
   - JWT_SECRET_KEY_ACCESS, JWT_SECRET_KEY_REFRESH
   - ADMIN, ADMIN_EMAIL, ADMIN_PASSWORD
4. Click each one to verify value
5. Check "Production" checkbox is selected

**Fix:**
1. Add missing variables
2. For existing: Edit and save again
3. Redeploy using "Redeploy" button

---

### Scenario 7: Can't Connect to Database From Vercel But Can Locally

**Symptoms:**
- Works on `localhost:8080`
- Fails when deployed to Vercel
- "connection refused" or "timed out"

**Likely Causes:**

| Cause | Solution |
|-------|----------|
| Using localhost as host | Change DB_HOST to public IP address |
| Database only accepts local connections | Configure database to accept remote |
| Firewall blocking Vercel | Add Vercel IP to firewall whitelist |
| Database provider requires IP whitelist | Add "%" or "0.0.0.0/0" to allow all |

**How to Fix:**

1. **Verify DB_HOST is NOT localhost**
   - Wrong: `DB_HOST=localhost` ❌
   - Right: `DB_HOST=db.example.com` ✅

2. **Check Database Firewall**
   - Neon: Settings → IP Whitelist → Add "0.0.0.0/0"
   - Railway: Doesn't require setup
   - Supabase: Database → Settings → SSL enforcement
   - AWS RDS: Security Groups → Inbound Rules → Add 0.0.0.0/0 for port 5432

3. **Test Locally First**
   ```bash
   # Test connection locally from terminal
   psql postgresql://user:pass@host:5432/dbname
   
   # If this works, issue is Vercel-specific
   # If this fails, it's a database issue
   ```

---

## ❓ Frequently Asked Questions

### Q: Is there a cost?
**A:** Free tier includes 100GB-hours per month of serverless execution. Most hobby projects use <10GB-hours/month. Database costs depend on provider (many have free tier).

### Q: Can I use a custom domain?
**A:** Yes! Settings → Domains → Add custom domain → Vercel shows DNS instructions.

### Q: How do I rollback to a previous version?
**A:** Deployments → Find previous version → Click ... menu → Redeploy.

### Q: Does auto-deployment work?
**A:** Yes, if you set up GitHub Actions. Every push to main branch auto-deploys.

### Q: How do I see what's happening?
**A:** Deployments → Click deployment → Logs tab shows real-time logs.

### Q: Can I monitor my API?
**A:** Yes! Click project → Analytics tab → See request counts, errors, response times.

### Q: What if I commit `.env` to git?
**A:** Don't worry, it's in `.gitignore`. But check with: `git ls-files | grep .env`.

### Q: Can I modify code on Vercel?
**A:** No, all changes must be made locally, committed, and pushed to git.

### Q: How do I update my API?
**A:** Make changes → `git commit` → `git push` → Vercel auto-deploys.

### Q: Is my data secure?
**A:** Yes, if you use HTTPS (automatic), strong passwords, and keep secrets in environment variables.

### Q: Can I upgrade later?
**A:** Yes, anytime. Start with free tier, upgrade to Pro ($20/month) for custom domains, advanced analytics, etc.

### Q: What if deployment is slow?
**A:** Normal first-time deployments take 3-5 minutes. Check Deployments → Logs for progress.

### Q: How often can I deploy?
**A:** Unlimited! Every git push = new deployment (usually within 1-2 minutes).

### Q: Can multiple people deploy?
**A:** Yes! All developers with git access can push and trigger deployments.

### Q: What about CORS?
**A:** Configure in your API handlers if needed. Currently uses default behavior.

---

## 🎯 Quick Fix Guide

**"Build failed"**
→ Check Deployments → Logs → Look for error → Fix locally → Push again

**"502 Bad Gateway"**
→ Check Logs for errors → Verify environment variables → Check database connection

**"Can't connect to database"**
→ Verify DB_HOST is not localhost → Check firewall → Test locally first

**"Environment variables not working"**
→ Settings → Environment Variables → Check all present → Verify Production selected → Redeploy

**"First request slow"**
→ Normal! (cold start) → Nothing to fix

**"Authentication failing"**
→ Check JWT secrets are set → Verify different from each other → Redeploy

**"Can't see logs"**
→ Deployments → Click latest → Logs tab (or Logs section)

**"Want to see old logs"**
→ Deployments → Click older deployment → Logs tab

**"Want to rollback"**
→ Deployments → Find old version → ... menu → Redeploy

---

## 📞 Support Resources

- **Vercel Docs**: https://vercel.com/docs
- **Go on Vercel**: https://vercel.com/docs/go/go-support
- **Vercel Status**: https://vercel-status.com
- **This Project Docs**:
  - Quick start: `VERCEL_QUICK_START.md`
  - Walkthrough: `DEPLOYMENT_WALKTHROUGH.md`
  - Detailed: `VERCEL_DEPLOYMENT_GUIDE.md`

---

## 🚨 Critical Checklist

Before contacting support:

- [ ] Checked Deployments → Logs for error messages
- [ ] Verified all environment variables are set
- [ ] Confirmed database credentials are correct
- [ ] Tested database connection locally
- [ ] Committed and pushed all changes
- [ ] Redeployed after making changes
- [ ] Waited at least 5 minutes after deployment
- [ ] Cleared browser cache (hard refresh: Ctrl+Shift+R)

---

**Most issues are resolved by:**
1. Checking logs
2. Verifying environment variables
3. Redeploying after changes

**Good luck!** 🚀
