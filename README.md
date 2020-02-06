![Logo](https://github.com/SysBind/chartpack/blob/master/logo.png)

# Air-Gapped Helm-Charts Packaging & Deployment

_Project Status:_ early alpha stage, in development

<pre>
  -----------           -----------
 |   Chart   | 1 --->* |   Image   |
  -----------           -----------
     *^
      |1
  ---------
 |  Repo   |
  ---------
</pre>

## Commands


### pack

Create a tarball with selected charts and images,
and also chartpack itself, for the deploy phase.


### deploy 

Copy images to selected targets and install charts 
