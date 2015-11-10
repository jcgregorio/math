Ideas
=====

The naming strategy for each element will be chapter-section-bucket.

Common code for the elements will be broken out into one or more
behaviors, such as practice, events, and animations.

The page will load with the correct elements linked, and
be given a series of element names to display, maybe a total
of 10. Then each will be displayed and then the results sent
back.

Move rint to common lib.
Move submit and continue to behavior.
Each element needs to implement isCorrect().
Does 1.2 implement slots? If not we need to add impl for that.

    <dom-module id="e-1-1-1">
      <template>
        <p>y={{m}}x+{{b}}</p>
        <p>What is the slope <number id=m></number> and the y-intercept <number id=b></number></p>
        <correction>
          The slope is {{m}} and the y-intercept is {{b}}.
        </correction>
      </template>
    </dom-module>
    <script>
      Polymer({
        is: "e-1-1-1",

        ready: function() {
          this.b = this._rint(2, 9);
          this.m = this._rint(2, 9);
        },

        _isCorrect: function() {
          return +this.$.m.value == this.m && +this.$.b.value == this.b;
        },
      });
    </script>

Need <correction> and <number> custom elements.
Need a common lib with rint() and near().
Need a behavior with common flows.
  - submit
  - continue
  - Injects the Submit and Continue buttons.

y=mx+b
y=mx-b
y=x/m`+b
y=mx
y=b
